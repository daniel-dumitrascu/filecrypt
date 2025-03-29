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

# Define the output binary names
CLIENT_BIN := client$(EXE)
SERVER_BIN := server$(EXE)

# Define the executables with their paths
CLIENT := $(CLIENT_DIR)$(SEPARATOR)$(CLIENT_BIN)
SERVER := $(SERVER_DIR)$(SEPARATOR)$(SERVER_BIN)

# Default target to build all executables
all: $(CLIENT) $(SERVER)

# Rule to build the client
$(CLIENT): $(wildcard $(CLIENT_DIR)/*.go)
	cd $(CLIENT_DIR) && $(GO) build -o $(CLIENT_BIN) .

# Rule to build the server
$(SERVER): $(wildcard $(SERVER_DIR)/*.go)
	cd $(SERVER_DIR) && $(GO) build -o $(SERVER_BIN) .
	
# Clean up the executables
clean:
	$(RMCMD) $(CLIENT) $(SERVER)