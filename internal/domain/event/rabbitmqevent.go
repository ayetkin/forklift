package event

type ExportVmMessage struct {
	MessageId string
}

type ConvertVmMessage struct {
	MessageId string
}

type CreateImageMessage struct {
	MessageId string
}

type CreateVolumeMessage struct {
	MessageId string
}

type CreateInstanceMessage struct {
	MessageId string
}

type ReserveFloatingIPMessage struct {
	MessageId string
}

type AssociateFloatingIPMessage struct {
	MessageId string
}

type CleanVmMessage struct {
	VmName string
}
