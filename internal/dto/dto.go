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
	Category   CategoryDTO `json:"category"`
	Answers    []AnswerDTO `json:"answers"`
}
