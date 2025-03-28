# Detect the operating system
ifeq ($(OS),Windows_NT)
    SEPARATOR := \\
    EXE := .exe
    RMCMD := del
else
    SEPARATOR := /
    EXE :=
    RMCMD := rm -f
endif

# Define the Go command
GO := go

# Define the directories
CLIENT_DIR := client
SERVER_DIR := server
TOOL_DIR := crypt

# Define the output binary names
CLIENT_BIN := client$(EXE)
SERVER_BIN := server$(EXE)
TOOL_BIN := crypt$(EXE)

# Define the executables with their paths
CLIENT := $(CLIENT_DIR)$(SEPARATOR)$(CLIENT_BIN)
SERVER := $(SERVER_DIR)$(SEPARATOR)$(SERVER_BIN)
CRYPT := $(TOOL_DIR)$(SEPARATOR)$(TOOL_BIN)

# Default target to build all executables
all: $(CLIENT) $(SERVER)

# Rule to build the client
$(CLIENT): $(wildcard $(CLIENT_DIR)/*.go)
	cd $(CLIENT_DIR) && $(GO) build -o $(CLIENT_BIN) .

# Rule to build the server
$(SERVER): $(wildcard $(SERVER_DIR)/*.go)
	cd $(SERVER_DIR) && $(GO) build -o $(SERVER_BIN) .

# Build the command line encrypt tool
tool: $(CRYPT)
$(CRYPT): $(wildcard $(TOOL_DIR)/*.go)
	cd $(TOOL_DIR) && $(GO) build -o $(TOOL_BIN) .
	
# Clean up the executables
clean:
	$(RMCMD) $(CRYPT) $(CLIENT) $(SERVER)