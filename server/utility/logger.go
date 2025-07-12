package utils

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

func loggerStyles() *log.Styles {
	styles := log.DefaultStyles()

	styles.Timestamp = lipgloss.NewStyle().Faint(true)
	styles.Key = lipgloss.NewStyle().Faint(true)
	styles.Value = lipgloss.NewStyle().Bold(true).Faint(true)
	styles.Levels = map[log.Level]lipgloss.Style{
		log.DebugLevel: lipgloss.NewStyle().
			SetString(" d ").
			Bold(true).MaxWidth(3).
			Background(lipgloss.Color("4")).
			Foreground(lipgloss.Color("8")),
		log.InfoLevel: lipgloss.NewStyle().
			SetString(" i ").
			Bold(true).MaxWidth(3).
			Background(lipgloss.Color("5")).
			Foreground(lipgloss.Color("8")),
		log.WarnLevel: lipgloss.NewStyle().
			SetString(" w ").
			Bold(true).MaxWidth(3).
			Background(lipgloss.Color("11")).
			Foreground(lipgloss.Color("8")),
		log.ErrorLevel: lipgloss.NewStyle().
			SetString(" e ").
			Bold(true).MaxWidth(3).
			Background(lipgloss.Color("9")).
			Foreground(lipgloss.Color("8")),
		log.FatalLevel: lipgloss.NewStyle().
			SetString(" f ").
			Bold(true).MaxWidth(3).
			Background(lipgloss.Color("1")).
			Foreground(lipgloss.Color("8")),
	}
	return styles
}

func SetupLogger() {
	log.SetStyles(loggerStyles())
	log.SetLevel(log.DebugLevel)

	if Config.Production {
		log.SetFormatter(log.JSONFormatter)
	} else {
		log.SetFormatter(log.TextFormatter)
	}
}
