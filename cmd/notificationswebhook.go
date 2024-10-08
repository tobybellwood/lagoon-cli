package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/uselagoon/machinery/api/lagoon"
	lclient "github.com/uselagoon/machinery/api/lagoon/client"
	"github.com/uselagoon/machinery/api/schema"

	"github.com/spf13/cobra"
	"github.com/uselagoon/lagoon-cli/pkg/output"
)

var addNotificationWebhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "Add a new webhook notification",
	Long: `Add a new webhook notification
This command is used to set up a new webhook notification in Lagoon. This requires information to talk to the webhook like the webhook URL.
It does not configure a project to send notifications to webhook though, you need to use project-webhook for that.`,
	Aliases: []string{"w"},
	PreRunE: func(_ *cobra.Command, _ []string) error {
		return validateTokenE(lagoonCLIConfig.Current)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		webhook, err := cmd.Flags().GetString("webhook")
		if err != nil {
			return err
		}
		organizationID, err := cmd.Flags().GetUint("organization-id")
		if err != nil {
			return err
		}
		if err := requiredInputCheck("Notification name", name, "Webhook", webhook); err != nil {
			return err
		}
		if yesNo(fmt.Sprintf("You are attempting to create a webhook notification '%s' with webhook url '%s', are you sure?", name, webhook)) {
			current := lagoonCLIConfig.Current
			token := lagoonCLIConfig.Lagoons[current].Token
			lc := lclient.New(
				lagoonCLIConfig.Lagoons[current].GraphQL,
				lagoonCLIVersion,
				lagoonCLIConfig.Lagoons[current].Version,
				&token,
				debug)

			notification := schema.AddNotificationWebhookInput{
				Name:         name,
				Webhook:      webhook,
				Organization: &organizationID,
			}

			result, err := lagoon.AddNotificationWebhook(context.TODO(), &notification, lc)
			if err != nil {
				return err
			}
			var data []output.Data
			notificationData := []string{
				returnNonEmptyString(fmt.Sprintf("%v", result.ID)),
				returnNonEmptyString(fmt.Sprintf("%v", result.Name)),
				returnNonEmptyString(fmt.Sprintf("%v", result.Webhook)),
			}
			if result.Organization != nil {
				organization, err := lagoon.GetOrganizationByID(context.TODO(), organizationID, lc)
				if err != nil {
					return err
				}
				notificationData = append(notificationData, organization.Name)
			} else {
				notificationData = append(notificationData, "-")
			}
			data = append(data, notificationData)
			r := output.RenderOutput(output.Table{
				Header: []string{
					"ID",
					"Name",
					"Webhook",
					"Organization",
				},
				Data: data,
			}, outputOptions)
			fmt.Fprintf(cmd.OutOrStdout(), "%s", r)
		}
		return nil
	},
}

var addProjectNotificationWebhookCmd = &cobra.Command{
	Use:     "project-webhook",
	Aliases: []string{"pw"},
	Short:   "Add a webhook notification to a project",
	Long: `Add a webhook notification to a project
This command is used to add an existing webhook notification in Lagoon to a project.`,
	PreRunE: func(_ *cobra.Command, _ []string) error {
		return validateTokenE(lagoonCLIConfig.Current)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		if err := requiredInputCheck("Notification name", name, "Project name", cmdProjectName); err != nil {
			return err
		}
		if yesNo(fmt.Sprintf("You are attempting to add webhook notification '%s' to project '%s', are you sure?", name, cmdProjectName)) {
			current := lagoonCLIConfig.Current
			token := lagoonCLIConfig.Lagoons[current].Token
			lc := lclient.New(
				lagoonCLIConfig.Lagoons[current].GraphQL,
				lagoonCLIVersion,
				lagoonCLIConfig.Lagoons[current].Version,
				&token,
				debug)
			notification := &schema.AddNotificationToProjectInput{
				NotificationType: schema.WebhookNotification,
				NotificationName: name,
				Project:          cmdProjectName,
			}
			_, err := lagoon.AddNotificationToProject(context.TODO(), notification, lc)
			if err != nil {
				return err
			}
			resultData := output.Result{
				Result: "success",
			}
			r := output.RenderResult(resultData, outputOptions)
			fmt.Fprintf(cmd.OutOrStdout(), "%s", r)
		}
		return nil
	},
}

var listProjectWebhooksCmd = &cobra.Command{
	Use:     "project-webhook",
	Aliases: []string{"pw"},
	Short:   "List webhook details about a project (alias: pw)",
	PreRunE: func(_ *cobra.Command, _ []string) error {
		return validateTokenE(lagoonCLIConfig.Current)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}
		if err := requiredInputCheck("Project name", cmdProjectName); err != nil {
			return err
		}

		current := lagoonCLIConfig.Current
		token := lagoonCLIConfig.Lagoons[current].Token
		lc := lclient.New(
			lagoonCLIConfig.Lagoons[current].GraphQL,
			lagoonCLIVersion,
			lagoonCLIConfig.Lagoons[current].Version,
			&token,
			debug)

		result, err := lagoon.GetProjectNotificationWebhook(context.TODO(), cmdProjectName, lc)
		if err != nil {
			return err
		}
		if len(result.Name) == 0 {
			return handleNilResults("No project found for '%s'\n", cmd, cmdProjectName)
		} else if len(result.Notifications.Webhook) == 0 {
			return handleNilResults("No webhook notificatons found for project: '%s'\n", cmd, cmdProjectName)
		}

		data := []output.Data{}
		if result.Notifications != nil {
			for _, notification := range result.Notifications.Webhook {
				data = append(data, []string{
					returnNonEmptyString(fmt.Sprintf("%v", notification.Name)),
					returnNonEmptyString(fmt.Sprintf("%v", notification.Webhook)),
				})
			}
		}
		r := output.RenderOutput(output.Table{
			Header: []string{
				"Name",
				"Webhook",
			},
			Data: data,
		}, outputOptions)
		fmt.Fprintf(cmd.OutOrStdout(), "%s", r)
		return nil
	},
}

var listAllWebhooksCmd = &cobra.Command{
	Use:     "webhook",
	Aliases: []string{"w"},
	Short:   "List all webhook notification details (alias: w)",
	PreRunE: func(_ *cobra.Command, _ []string) error {
		return validateTokenE(lagoonCLIConfig.Current)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}
		current := lagoonCLIConfig.Current
		token := lagoonCLIConfig.Lagoons[current].Token
		lc := lclient.New(
			lagoonCLIConfig.Lagoons[current].GraphQL,
			lagoonCLIVersion,
			lagoonCLIConfig.Lagoons[current].Version,
			&token,
			debug)
		result, err := lagoon.GetAllNotificationWebhook(context.TODO(), lc)
		if err != nil {
			return err
		}
		data := []output.Data{}
		for _, res := range *result {
			b, _ := json.Marshal(res.Notifications.Webhook)
			if string(b) != "null" {
				for _, notif := range res.Notifications.Webhook {
					data = append(data, []string{
						returnNonEmptyString(fmt.Sprintf("%v", res.Name)),
						returnNonEmptyString(fmt.Sprintf("%v", notif.Name)),
						returnNonEmptyString(fmt.Sprintf("%v", notif.Webhook)),
					})
				}
			}
		}
		r := output.RenderOutput(output.Table{
			Header: []string{
				"Project",
				"Name",
				"Webhook",
			},
			Data: data,
		}, outputOptions)
		fmt.Fprintf(cmd.OutOrStdout(), "%s", r)
		return nil
	},
}

var deleteProjectWebhookNotificationCmd = &cobra.Command{
	Use:     "project-webhook",
	Aliases: []string{"pw"},
	Short:   "Delete a webhook notification from a project",
	PreRunE: func(_ *cobra.Command, _ []string) error {
		return validateTokenE(lagoonCLIConfig.Current)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		if err := requiredInputCheck("Project name", cmdProjectName, "Notification name", name); err != nil {
			return err
		}

		current := lagoonCLIConfig.Current
		token := lagoonCLIConfig.Lagoons[current].Token
		lc := lclient.New(
			lagoonCLIConfig.Lagoons[current].GraphQL,
			lagoonCLIVersion,
			lagoonCLIConfig.Lagoons[current].Version,
			&token,
			debug)

		project, err := lagoon.GetProjectByName(context.TODO(), cmdProjectName, lc)
		if err != nil {
			return err
		}
		if project.Name == "" {
			return handleNilResults("No project found for '%s'\n", cmd, cmdProjectName)
		}

		if yesNo(fmt.Sprintf("You are attempting to delete webhook notification '%s' from project '%s', are you sure?", name, cmdProjectName)) {
			notification := &schema.RemoveNotificationFromProjectInput{
				NotificationType: schema.WebhookNotification,
				NotificationName: name,
				Project:          cmdProjectName,
			}
			_, err := lagoon.RemoveNotificationFromProject(context.TODO(), notification, lc)
			if err != nil {
				return err
			}
			resultData := output.Result{
				Result: "success",
			}
			r := output.RenderResult(resultData, outputOptions)
			fmt.Fprintf(cmd.OutOrStdout(), "%s", r)
		}
		return nil
	},
}

var deleteWebhookNotificationCmd = &cobra.Command{
	Use:     "webhook",
	Aliases: []string{"w"},
	Short:   "Delete a webhook notification from Lagoon",
	PreRunE: func(_ *cobra.Command, _ []string) error {
		return validateTokenE(lagoonCLIConfig.Current)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		if err := requiredInputCheck("Notification name", name); err != nil {
			return err
		}
		// Todo: Verify notifcation name exists - requires #PR https://github.com/uselagoon/lagoon/pull/3740
		if yesNo(fmt.Sprintf("You are attempting to delete webhook notification '%s', are you sure?", name)) {
			current := lagoonCLIConfig.Current
			token := lagoonCLIConfig.Lagoons[current].Token
			lc := lclient.New(
				lagoonCLIConfig.Lagoons[current].GraphQL,
				lagoonCLIVersion,
				lagoonCLIConfig.Lagoons[current].Version,
				&token,
				debug)
			result, err := lagoon.DeleteNotificationWebhook(context.TODO(), name, lc)
			if err != nil {
				return err
			}
			resultData := output.Result{
				Result: result.DeleteNotification,
			}
			r := output.RenderResult(resultData, outputOptions)
			fmt.Fprintf(cmd.OutOrStdout(), "%s", r)
		}
		return nil
	},
}

var updateWebhookNotificationCmd = &cobra.Command{
	Use:     "webhook",
	Aliases: []string{"w"},
	Short:   "Update an existing webhook notification",
	PreRunE: func(_ *cobra.Command, _ []string) error {
		return validateTokenE(lagoonCLIConfig.Current)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			return err
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		newname, err := cmd.Flags().GetString("newname")
		if err != nil {
			return err
		}
		webhook, err := cmd.Flags().GetString("webhook")
		if err != nil {
			return err
		}
		if err := requiredInputCheck("Notification name", name); err != nil {
			return err
		}
		patch := schema.UpdateNotificationWebhookPatchInput{
			Name:    nullStrCheck(newname),
			Webhook: nullStrCheck(webhook),
		}
		if patch == (schema.UpdateNotificationWebhookPatchInput{}) {
			return fmt.Errorf("missing arguments: either webhook or newname must be defined")
		}

		if yesNo(fmt.Sprintf("You are attempting to update webhook notification '%s', are you sure?", name)) {
			current := lagoonCLIConfig.Current
			token := lagoonCLIConfig.Lagoons[current].Token
			lc := lclient.New(
				lagoonCLIConfig.Lagoons[current].GraphQL,
				lagoonCLIVersion,
				lagoonCLIConfig.Lagoons[current].Version,
				&token,
				debug)

			notification := &schema.UpdateNotificationWebhookInput{
				Name:  name,
				Patch: patch,
			}
			result, err := lagoon.UpdateNotificationWebhook(context.TODO(), notification, lc)
			if err != nil {
				return err
			}
			data := []output.Data{
				[]string{
					returnNonEmptyString(fmt.Sprintf("%v", result.ID)),
					returnNonEmptyString(fmt.Sprintf("%v", result.Name)),
					returnNonEmptyString(fmt.Sprintf("%v", result.Webhook)),
				},
			}
			r := output.RenderOutput(output.Table{
				Header: []string{
					"ID",
					"Name",
					"Webhook",
				},
				Data: data,
			}, outputOptions)
			fmt.Fprintf(cmd.OutOrStdout(), "%s", r)
		}
		return nil
	},
}

func init() {
	addNotificationWebhookCmd.Flags().StringP("name", "n", "", "The name of the notification")
	addNotificationWebhookCmd.Flags().StringP("webhook", "w", "", "The webhook URL of the notification")
	addNotificationWebhookCmd.Flags().Uint("organization-id", 0, "ID of the Organization")
	addProjectNotificationWebhookCmd.Flags().StringP("name", "n", "", "The name of the notification")
	deleteProjectWebhookNotificationCmd.Flags().StringP("name", "n", "", "The name of the notification")
	deleteWebhookNotificationCmd.Flags().StringP("name", "n", "", "The name of the notification")
	updateWebhookNotificationCmd.Flags().StringP("name", "n", "", "The name of the notification")
	updateWebhookNotificationCmd.Flags().StringP("newname", "N", "", "The name of the notification")
	updateWebhookNotificationCmd.Flags().StringP("webhook", "w", "", "The webhook URL of the notification")
}
