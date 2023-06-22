package netkit

var (
	RDWR_BUFFER_SIZE     = 10 << 10
	Max_Error_Time       = 100
	Single_Max_Err_Times = 10
)

type DataContent struct {
	Data  []byte
	Types string
}
