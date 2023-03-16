package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	crdb "github.com/labring/crdbase"
	"github.com/labring/crdbase/tests/examples/models"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "test.sealos.io", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

var (
	mainScheme = runtime.NewScheme()
	setupLog   = ctrl.Log.WithName("setup")
)

func initManager() ctrl.Manager {
	utilruntime.Must(clientgoscheme.AddToScheme(mainScheme))
	utilruntime.Must(AddToScheme(mainScheme))

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:         mainScheme,
		Port:           9443,
		LeaderElection: false,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	return mgr
}

func main() {
	mgr := initManager()

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	go func() {
		setupLog.V(1).Info("starting manager")
		if err := mgr.Start(ctx); err != nil {
			setupLog.Error(err, "problem running manager")
			os.Exit(1)
		}
	}()

	done := make(chan struct{})
	go func() {
		if mgr.GetCache().WaitForCacheSync(context.Background()) {
			done <- struct{}{}
		}
	}()
	<-done

	// node := &apiextv1.CustomResourceDefinitionList{}
	// if err := mgr.GetClient().List(ctx, node); err != nil {
	// 	panic(err)
	// }
	// fmt.Println(node.Items)

	conf := crdb.CRDBaseConfig{
		Manager: mgr,
		GroupVersion: schema.GroupVersion{
			Group:   "test.sealos.io",
			Version: "v1",
		},
		ServiceAccount: "crdb-test",
		Namespace:      "crdb-test",
	}

	db, err := crdb.NewCRDBase(conf, setupLog)
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(ctx, models.Count{}); err != nil {
		setupLog.V(1).Info("unable to auto migrate", "error", err)
	}

	dbCount := &models.Count{}

	got, err := dbCount.Add(ctx, db, "user1", models.DownloadCount, 1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(got)
}
