INIT_SCRIPT := ./scripts/!init_project.sh
BUILD_SCRIPT := ./scripts/build.sh
RUN_SCRIPT := ./scripts/run.sh

init:
	@if [ ! -x "$(INIT_SCRIPT)" ]; then\
		chmod +x "$(INIT_SCRIPT)";\
	fi
	"$(INIT_SCRIPT)"
build:
	@if [ ! -x "$(BUILD_SCRIPT)" ]; then\
		chmod +x "$(BUILD_SCRIPT)";\
	fi
	"$(BUILD_SCRIPT)"
run:
	@if [ ! -x "$(RUN_SCRIPT)" ]; then\
		chmod +x "$(RUN_SCRIPT)";\
	fi
	"$(RUN_SCRIPT)"
