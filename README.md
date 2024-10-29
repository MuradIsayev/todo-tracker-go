
# Todo Tracker CLI in Go

Welcome to todo-tracker-go! This project is a CLI-based task management tool written in Go. Designed to help users manage tasks and projects with ease, it provides a minimalistic, terminal-based interface for organizing tasks, tracking progress, and even setting timers. This tool is perfect for developers who need a lightweight, yet powerful tracker for their coding projects.


## Screenshots

_*Projects Table*_

![Projects Table](https://i.postimg.cc/ZRzmvKCL/image.png "Projects Table")
&nbsp;

_*Tasks Table with controlled timer functionality*_

![Tasks Table with Coutdown](https://i.postimg.cc/KYtTMfZ3/image.png "Tasks Table with Countdown")


## Features

* **Project Management:** Create, update, delete, list and manage projects, allowing for organized oversight of your work.

- **Task Tracking:** Add tasks to projects, update task details, and delete tasks within the REPL mode, maintaining full control over your tasks.

- **Timer Functionality:** Set countdown timers for tasks and manage timer states with commands to start, stop, pause, and resume within the Timer mode.

- **Persistent Storage:** Data is stored in a JSON file, ensuring all projects and tasks persist between sessions.

- **Modular Codebase:** The code is organized with interfaces to reduce repetition, particularly between the task and project modules.

- **User-Friendly Help Command:** Access a comprehensive help guide using the help command, providing detailed descriptions of commands and examples for beginners.

## Run Locally

Clone the project and go to the project directory:

```bash
git clone https://github.com/MuradIsayev/todo-tracker-go.git
cd todo-tracker-go

```

Build the project:

```bash
  go build -o todo-tracker
```

Run the executable:

```bash
  ./todo-tracker
```


