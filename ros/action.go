package ros

type ActionType interface {
	MD5Sum() string
	Name() string
	GoalType() MessageType
	FeedbackType() MessageType
	ResultType() MessageType
	NewAction() Action
}

type Action interface {
	GoalMessage() Message
	FeedbackMessage() Message
	ResultMessage() Message
}
