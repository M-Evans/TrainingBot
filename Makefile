all:
	go build

install:
	sudo cp trainingbot.service /etc/systemd/system/trainingbot.service

