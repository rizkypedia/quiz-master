package dto

type QuizTypeDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CategoryDTO struct {
	ID       uint        `json:"id"`
	Name     string      `json:"name"`
	QuizType QuizTypeDTO `json:"quizType"`
}

type AnswerDTO struct {
	ID         uint   `json:"id"`
	Answer     string `json:"answer"`
	QuestionId uint   `json:"questionId"`
	IsCorrect  bool   `json:"isCorrect"`
}

type QuestionDTO struct {
	ID         uint        `json:"id"`
	Question   string      `json:"question"`
	CategoryId uint        `json:"categoryId"`
	Category   CategoryDTO `json:"category,omitempty"`
	Answers    []AnswerDTO `json:"answers,omitempty"`
}

type CreateCategoryDTO struct {
	ID       uint         `json:"id"`
	Name     string       `json:"name"`
	QuizType *QuizTypeDTO `json:"quizType,omitempty"`
}

type CreateAnswerDTO struct {
	Answer    *string `json:"answer" binding:"required"`
	IsCorrect *bool   `json:"isCorrect" binding:"required"`
}

type CreateQuestionDTO struct {
	Question   *string            `json:"question" binding:"required"`
	CategoryId *uint              `json:"categoryId" binding:"required"`
	Answers    *[]CreateAnswerDTO `json:"answers" binding:"required"`
}
