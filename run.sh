#!/bin/sh

printUsage(){
	echo "printUsage"
	echo " (+)startall:    start all your APP"
	echo " (+)stopall:     stop all your APP"
	echo " (+)build:       build all your APP"
	echo " (+)genconfig:   generate your serverconfig from your config/devconfig.json"
}

case "$1" in
    -h|-?|h|help)
		printUsage
		;;
    build)
        cd example
        make build
		;;
	startall)
        cd example
        sh run.sh start
		;;
	stopall)
        cd example
        sh run.sh stop
		;;
	genconfig)
        # build gen
        export GOPATH="$(PWD)/.."
        go build -o genconfig "./tool/genconfig/"
	    app="./genconfig"
	    configDir="./config/devconfig.json"
	    if [ $# -lt 2 ]; then
		    ${app} -config="${configDir}"
		else
		    ${app} -config="${configDir}" -target="$2"
	    fi
	    rm ${app}
		;;
	*)
		echo "unknown Operation: $1"
		printUsage
		;;
esac

exit 0
