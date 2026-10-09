package dto

type QuizTypeDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CategoryDTO struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type AnswerDTO struct {
	ID         uint   `json:"id"`
	Answer     string `json:"answer"`
	QuestionId uint   `json:"questionId"`
}

type QuestionDTO struct {
	ID         uint        `json:"id"`
	Question   string      `json:"question"`
	CategoryId uint        `json:"categoryId"`
	Category   CategoryDTO `json:"category,omitempty"`
	QuizTypeID uint        `json:"quizTypeId"`
	QuizType   QuizTypeDTO `json:"quizType"`
	Answers    []AnswerDTO `json:"answers"`
}

type QuizAnswerSubmission struct {
	QuestionID uint   `json:"questionId" binding:"required"`
	AnswerIDs  []uint `json:"answerIds"`
}

type QuizResultDTO struct {
	CategoryID   uint `json:"categoryId"`
	CorrectCount int  `json:"correctCount"`
	TotalCount   int  `json:"totalCount"`
	Percentage   int  `json:"percentage"`
}

type CreateCategoryDTO struct {
	Name string `json:"name" binding:"required"`
}

type SetCategoryActiveDTO struct {
	Active *bool `json:"active" binding:"required"`
}

type CreateAnswerDTO struct {
	Answer    *string `json:"answer" binding:"required"`
	IsCorrect *bool   `json:"isCorrect" binding:"required"`
}

type CreateQuestionDTO struct {
	Question   *string            `json:"question" binding:"required"`
	CategoryId *uint              `json:"categoryId" binding:"required"`
	QuizTypeID *uint              `json:"quizTypeId" binding:"required"`
	Answers    *[]CreateAnswerDTO `json:"answers" binding:"required"`
}

type AuthRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
