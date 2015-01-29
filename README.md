godbtoswift
============================

Simple utility that checks for the latest file and pushes it to an OpenStack ObjectStore.

Wrote it for fun and to easily backup my database files to an OpenStack ObjectStore.

Usage:

	godbtoswift -s=/backups/daily/ -t=mydbname -c=/path/to/config.json

	-s=/backups/daily/ -- Place where the application can find the backups
	-t=mydbname -- Directory it will place the files in, at your ObjectStore
	-c=/path/to/config.json -- Location of config file
