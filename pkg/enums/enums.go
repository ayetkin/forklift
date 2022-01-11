package enums

type Status string

const (
	Pending   Status = "Migration Task Pending Queue"
	Running          = "Migration Task Running"
	Completed        = "Migration Task Completed"
	Failed           = "Migration Task Filed"
	Deleted          = "Migration Task Deleted"
)

type Stage string

const (
	PendingQueue        Stage = "Migration Pending Queue"
	ExportVm                  = "Exporting VM"
	ConvertVm                 = "Converting VM"
	CreateImage               = "Creating Image"
	CreateVolume              = "Creating Volume"
	CreateInstance            = "Creating Instance"
	ReserveFloatingIP         = "Reserving Floating IP"
	AssociateFloatingIP       = "Associating Floating IP"
	Finished                  = "Migration Finished"
)
