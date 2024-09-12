package test

import "testing"
import "github.com/Pooh-Mucho/go-aigc/chat"

func Test_DashScope_QwenMax_Hello(t *testing.T) {
	Tests.Hello(t, chat.Models.QwenMax_20240428, WithAliyun)
}

func Test_DashScope_QwenMax_Chinese_Poetry(t *testing.T) {
	Tests.Chinese_Poetry(t, chat.Models.QwenMax_20240428, WithAliyun)
}

func Test_DashScope_QwenMax_Emoji(t *testing.T) {
	Tests.Emoji(t, chat.Models.QwenMax_20240428, WithAliyun)
}

func Test_DashScope_QwenMax_Multi_Contents(t *testing.T) {
	Tests.Multi_Contents(t, chat.Models.QwenMax_20240428, WithAliyun)
}

func Test_DashScope_QwenMax_Template_System_Injections(t *testing.T) {
	for i := 0; i < 1; i++ {
		Tests.Template_System_Injections(t, chat.Models.QwenMax_20240428, WithAliyun)
	}
}

func Test_DashScope_QwenMax_System_Injections(t *testing.T) {
	for i := 0; i < 1; i++ {
		Tests.Template_System_Injections(t, chat.Models.QwenMax_20240428, WithAliyun)
	}
}

func Test_DashScope_QwenMax_Tool_Random_Number(t *testing.T) {
	Tests.Tool_Random_Number(t, chat.Models.QwenMax_20240428, WithAliyun)
}

func Test_DashScope_QwenMax_Tool_Add_Single(t *testing.T) {
	Tests.Tool_Add_Single(t, chat.Models.QwenMax_20240428, WithAliyun)
}

func Test_DashScope_QwenMax_Tool_Add_Parallel(t *testing.T) {
	Tests.Tool_Add_Parallel(t, chat.Models.QwenMax_20240428, WithAliyun)
}

func Test_DashScope_QwenMax_Tool_Sum(t *testing.T) {
	Tests.Tool_Sum(t, chat.Models.QwenMax_20240428, WithAliyun)
}

func Test_DashScope_QwenMax_Tool_File(t *testing.T) {
	Tests.Tool_File(t, chat.Models.QwenMax_20240428, WithAliyun)
}
