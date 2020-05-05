package tor
import(
	"github.com/digantaTheProgrammer/NodeTor/src/nodejs/apt"
	"github.com/cloudfoundry/libbuildpack"
)

func InstallTor(s *apt.Supplier,Log *libbuildpack.Logger) error {
	Log.Info("Installing Tor.....")
	err:=apt.SingleInstall(s,"tor","repo")
	if(err!=nil){
		return err
	}
	Log.Info("Tor installed!!!")
		
	torscript:=`export TOR_PORT_1=58`
	s.Stager.WriteProfileD("tor.sh",torscript)
	return nil
}