import api

# set to false to let the game run until ended
WAIT_FOR_KEY = True

# start new game
id = api.new_game()["id"]
print(f"game: {id}")

while True:
    game_state = api.get_state(id)
    state_name = game_state["state_name"]
    active_player = game_state["active_player"]
    has_ended = game_state["has_ended"]
    if has_ended:
        break

    print(f"State: {state_name} (player {active_player} to move)")

    if WAIT_FOR_KEY:
        input("Press enter to execute next move")

    # pick the first action for player 1
    api.send_actions(id, {f"{active_player}": 0})

    if WAIT_FOR_KEY:
        input("Press enter to execute next move")

    # pick the first actions for other players
    api.send_actions(id, {f"{p}": 0 for p in [1, 2, 3, 4] if p != active_player})
