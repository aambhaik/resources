package util

/**
Mashling gateway constants
*/
const Gateway_Definition_File_Name = "mashling.json"

const Gateway_Trigger_Config_Ref_Key = "config"
const Gateway_Trigger_Config_Prefix = "${configurations."
const Gateway_Trigger_Config_Suffix = "}"
const Gateway_Trigger_Handler_UseReplyHandler = "useReplyHandler"
const Gateway_Trigger_Handler_UseReplyHandler_Default = "false"
const Gateway_Trigger_Handler_AutoIdReply = "autoIdReply"
const Gateway_Trigger_Handler_AutoIdReply_Default = "false"
const Gateway_Trigger_Metadata_JSON_Name = "trigger.json"

const Gateway_Link_Condition_Operator_Equals = "=="
const Gateway_Link_Condition_Operator_NotEquals = "!="
const Gateway_Link_Condition_LHS_Start_Expr = "${"
const Gateway_Link_Condition_LHS_End_Expr = "}"
const Gateway_JSON_Content_Root_Env_Key = "TRIGGER_CONTENT_ROOT"
const Gateway_Link_Condition_LHS_JSON_Content_Prefix_Default = "trigger.content"
const Gateway_Link_Condition_LHS_JSONPath_Root = "$"

/**
Flogo constants
*/
const Flogo_App_Type = "flogo:app"
const Flogo_Trigger_Handler_Setting_Condition = "Condition"
