import api

# set to false to let the game run until ended
WAIT_FOR_KEY = True

# start new game
id = api.new_game()["id"]
print(f"game: {id}")

while True:
    game_state = api.get_state(id)
    state_name = game_state["state_name"]
    active_players = game_state["active_players"]
    has_ended = game_state["has_ended"]
    if has_ended:
        break

    print(f"State: {state_name} (player(s) {active_players} to move)")

    if WAIT_FOR_KEY:
        input("Press enter to execute next move")

    for p in active_players:
        actions = dict()

        # just pick the first action for now
        actions[f"{p}"] = 0
        api.send_actions(id, actions)
