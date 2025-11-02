INSTALL_SCRIPT := ./scripts/install_project.sh
BUILD_SCRIPT := ./scripts/build.sh

install:
	@if [ ! -x "$(INSTALL_SCRIPT)" ]; then\
		chmod +x "$(INSTALL_SCRIPT)";\
	fi
	"$(INSTALL_SCRIPT)"
build:
	@if [ ! -x "$(BUILD_SCRIPT)" ]; then\
		chmod +x "$(BUILD_SCRIPT)";\
	fi
	"$(BUILD_SCRIPT)"
