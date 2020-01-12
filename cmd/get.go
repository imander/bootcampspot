package cmd

import (
	"errors"

	"github.com/imander/bootcampspot/bcs"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get bootcamp spot details",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("no argument provided")
	},
}

var getAttendanceCmd = &cobra.Command{
	Use:     "attendance",
	Aliases: []string{"a"},
	Short:   "This command is used to get student attendance statistics",
	Run: func(cmd *cobra.Command, args []string) {
		setAll()
		attendance := bcs.GetAttendance()
		attendance.Metrics().Print()
	},
}

var getAssignmentsCmd = &cobra.Command{
	Use:     "assignments",
	Aliases: []string{"assignment", "as"},
	Short:   "This command is used to get student assignment submissions",
	Run: func(cmd *cobra.Command, args []string) {
		setEnrollmentID()
		assignments := bcs.GetGrades()
		assignments.Metrics().Print()
	},
}

var getFeedbackCmd = &cobra.Command{
	Use:     "feedback",
	Aliases: []string{"f"},
	Short:   "This command is used to get student feedback responses",
	Run: func(cmd *cobra.Command, args []string) {
		setCourseID()
		feedbacks := bcs.GetFeedback()
		feedbacks.Print()
	},
}

var getEnrollmentsCmd = &cobra.Command{
	Use:     "enrollments",
	Aliases: []string{"enrollment", "e"},
	Short:   "This command is used to get user course enrollments",
	Run: func(cmd *cobra.Command, args []string) {
		user := bcs.GetUser()
		user.PrintEnrollments()
	},
}

func setCourseID() {
	if bcs.CourseID != -1 {
		return
	}
	u := bcs.GetUser()
	u.ChooseEnrollment()
}

func setEnrollmentID() {
	if bcs.EnrollmentID != -1 {
		return
	}
	u := bcs.GetUser()
	u.ChooseEnrollment()
}

func setAll() {
	if bcs.CourseID != -1 && bcs.EnrollmentID != -1 {
		return
	}
	u := bcs.GetUser()
	u.ChooseEnrollment()
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(getAttendanceCmd)
	getCmd.AddCommand(getAssignmentsCmd)
	getCmd.AddCommand(getFeedbackCmd)
	getCmd.AddCommand(getEnrollmentsCmd)

	getFeedbackCmd.Flags().IntVarP(&bcs.CourseID, "course-id", "c", -1, "Course ID (optional)")
	getAttendanceCmd.Flags().IntVarP(&bcs.CourseID, "course-id", "c", -1, "Course ID (optional)")

	getAssignmentsCmd.Flags().IntVarP(&bcs.EnrollmentID, "enrollment-id", "e", -1, "Enrollment ID (optional)")
	getAttendanceCmd.Flags().IntVarP(&bcs.EnrollmentID, "enrollment-id", "e", -1, "Enrollment ID (optional)")
}
