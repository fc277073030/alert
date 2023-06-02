package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fc277073030/alert/models"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	webhookURL string
)

// prometheusCmd represents the prometheus command
var prometheusCmd = &cobra.Command{
	Use:   "prometheus",
	Short: "Prometheus alertmanager notification",
	Long:  "This is an alert from promtheus alertmanage",
	Run:   runPrometheusHandler,
}

func init() {
	rootCmd.AddCommand(prometheusCmd)

	prometheusCmd.Flags().StringVar(&webhookURL, "webhook", "", "Webhook URL for sending alerts")
	err := prometheusCmd.MarkFlagRequired("webhook")
	if err != nil {
		return
	}
}

func runPrometheusHandler(cmd *cobra.Command, _ []string) {
	_, _ = cmd.Flags().GetString("webhook")

	router := setupRouter()
	address := ":8080"

	log.Println("Server listening on", address)
	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/webhook", func(c *gin.Context) {
		handleAlert(c)
	})

	return router
}

func handleAlert(c *gin.Context) {
	// Handle alert data here
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}

	var alertData models.AlertData
	err = json.Unmarshal(body, &alertData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	formattedMessage := formatAlertMessage(alertData)
	sendAlertToWebhook(formattedMessage)
	c.JSON(http.StatusOK, gin.H{"message": "Alert processed successfully"})
}

func formatAlertMessage(alertData models.AlertData) string {
	alertsMessage := "告警列表：\n"
	for _, alert := range alertData.Alerts {
		alertMessage := fmt.Sprintf(
			"告警名称：%s\n"+
				"描述：%s\n"+
				"摘要：%s\n"+
				"开始时间：%s\n"+
				"恢复时间：%s\n\n",
			alert.Labels.Alertname,
			alert.Annotations.Description,
			alert.Annotations.Summary,
			alert.StartsAt,
			alert.EndsAt)
		alertsMessage += alertMessage
	}

	message := fmt.Sprintf(
		"状态：%s\n\n"+
			"告警名称：%s\n"+
			"描述：%s\n"+
			"摘要：%s\n"+
			"开始时间：%s\n"+
			"恢复时间：%s\n",
		alertData.Status,
		alertData.Alerts[0].Labels.Alertname,
		alertData.Alerts[0].Annotations.Description,
		alertData.Alerts[0].Annotations.Summary,
		alertData.Alerts[0].StartsAt,
		alertData.Alerts[0].EndsAt)

	return message
}

func sendAlertToWebhook(message string) {
	payload := map[string]interface{}{
		"text":       message,
		"markdown":   true,
		"username":   "Alert Bot",
		"icon_emoji": ":exclamation:",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Failed to marshal payload:", err)
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Failed to send message to Webhook:", err)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body:", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Webhook request failed with status:", resp.StatusCode)
		return
	}

	fmt.Println("Message sent to Webhook successfully")
}
