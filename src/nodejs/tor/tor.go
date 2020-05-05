package tor
import(
	"github.com/digantaTheProgrammer/NodeTor/src/nodejs/apt"
	"github.com/cloudfoundry/libbuildpack"
	"path/filepath"
)

func InstallTor(s *apt.Supplier,Log *libbuildpack.Logger) error {
	Log.Info("Installing Tor.....")
	apt.SingleInstall(s,"tor")
	Log.Info("Tor installed!!!")
	
	torscript:=`export TOR_PORT_1=58`
	s.Stager.WriteProfileD("tor.sh",torscript)
	return nil
}