package dto

type QuizTypeDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type CategoryDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
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
	QuizTypeID uint        `json:"quizTypeId"`
	QuizType   QuizTypeDTO `json:"quizType"`
	Answers    []AnswerDTO `json:"answers,omitempty"`
}

type CreateCategoryDTO struct {
	Name string `json:"name" binding:"required"`
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
