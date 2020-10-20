package errors

var ErrCanNotNil = New("shouldn't be nil")    //不能为空
var ErrMustPointer = New("Must be a pointer") //必须是指针
