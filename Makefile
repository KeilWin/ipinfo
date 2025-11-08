INIT_SCRIPT := ./scripts/!init_project.sh
BUILD_IPINFO_SCRIPT := ./scripts/build_ipinfo.sh
BUILD_IPINFO_UPDATER_SCRIPT := ./scripts/build_ipinfo_updater.sh
RUN_IPINFO_SCRIPT := ./scripts/run_ipinfo.sh
RUN_IPINFO_UPDATER_SCRIPT := ./scripts/run_ipinfo_updater.sh

init:
	@if [ ! -x "$(INIT_SCRIPT)" ]; then\
		chmod +x "$(INIT_SCRIPT)";\
	fi
	"$(INIT_SCRIPT)"
build-ipinfo:
	@if [ ! -x "$(BUILD_IPINFO_SCRIPT)" ]; then\
		chmod +x "$(BUILD_IPINFO_SCRIPT)";\
	fi
	"$(BUILD_IPINFO_SCRIPT)"
build-ipinfo-updater:
	@if [ ! -x "$(BUILD_IPINFO_UPDATER_SCRIPT)" ]; then\
		chmod +x "$(BUILD_IPINFO_UPDATER_SCRIPT)";\
	fi
	"$(BUILD_IPINFO_UPDATER_SCRIPT)"
run-ipinfo:
	@if [ ! -x "$(RUN_IPINFO_SCRIPT)" ]; then\
		chmod +x "$(RUN_IPINFO_SCRIPT)";\
	fi
	"$(RUN_IPINFO_SCRIPT)"
run-ipinfo-updater:
	@if [ ! -x "$(RUN_IPINFO_UPDATER_SCRIPT)" ]; then\
		chmod +x "$(RUN_IPINFO_UPDATER_SCRIPT)";\
	fi
	"$(RUN_IPINFO_UPDATER_SCRIPT)"
