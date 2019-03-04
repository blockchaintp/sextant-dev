#!/bin/bash
SESSION=sextant

function start() {
  if tmux has-session -t "$SESSION" 2>/dev/null; then
    echo "Session $SESSION already exists. Attaching..."
    sleep 1
    tmux -2 attach -t $SESSION
    exit 0;
  fi

  echo "Starting docker-compose"
  export MANUALRUN=1
  make dev

  echo "Creating tmux session $SESSION..."

  # get the size of the window and create a session at that size
  set -- $(stty size)
  tmux -2 new-session -d -s $SESSION -x "$2" -y "$(($1 - 1))"

  # the right hand col with a 50% vertical split
  tmux split-window -h -d
  tmux select-pane -t 1
  tmux split-window -v -d
  tmux select-pane -t 0

  # the left hand col with a row for each service
  tmux split-window -v -d

  tmux send-keys -t 0 'make frontend.run' C-m
  tmux send-keys -t 1 'make api.run' C-m
  
  tmux -2 attach-session -t $SESSION
}

function stop() {
  echo "Stopping tmux session $SESSION..."
  tmux kill-session -t $SESSION
  echo "Removing docker containers"
  docker rm -f $(docker ps -aq)
}

command="$@"

if [ -z "$command" ]; then
  command="start"
fi

eval "$command"