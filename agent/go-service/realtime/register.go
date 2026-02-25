package realtime

import "github.com/MaaXYZ/maa-framework-go/v4"

var (
	_ maa.CustomRecognitionRunner = &RealTimeAutoFightEntryRecognition{}
	_ maa.CustomRecognitionRunner = &RealTimeAutoFightExitRecognition{}
	_ maa.CustomRecognitionRunner = &RealTimeAutoFightSkillRecognition{}
	_ maa.CustomRecognitionRunner = &RealTimeAutoFightEndSkillRecognition{}
	_ maa.CustomActionRunner      = &RealTimeAutoFightSkillAction{}
	_ maa.CustomActionRunner      = &RealTimeAutoFightEndSkillAction{}
)

// 后续和AutoFight通用模块合并
func Register() {
	maa.AgentServerRegisterCustomRecognition("RealTimeAutoFightEntryRecognition", &RealTimeAutoFightEntryRecognition{})
	maa.AgentServerRegisterCustomRecognition("RealTimeAutoFightExitRecognition", &RealTimeAutoFightExitRecognition{})
	maa.AgentServerRegisterCustomRecognition("RealTimeAutoFightSkillRecognition", &RealTimeAutoFightSkillRecognition{})
	maa.AgentServerRegisterCustomAction("RealTimeAutoFightSkillAction", &RealTimeAutoFightSkillAction{})
	maa.AgentServerRegisterCustomRecognition("RealTimeAutoFightEndSkillRecognition", &RealTimeAutoFightEndSkillRecognition{})
	maa.AgentServerRegisterCustomAction("RealTimeAutoFightEndSkillAction", &RealTimeAutoFightEndSkillAction{})
}
