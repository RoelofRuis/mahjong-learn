import api
import strategy

# set to false to let the game run until ended
WAIT_FOR_KEY = True

# select the strategy to use, see strategy file for options
ACTION_SELECTION_STRATEGY = strategy.pick_first_action

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

    print(f" > State: {state_name} (player(s) {active_players} to move)")

    selected_action_ids = dict()
    selected_action_names = list()
    for p in active_players:
        allowed_actions = game_state["players"][str(p)]["actions"]
        selected_action_id = ACTION_SELECTION_STRATEGY(allowed_actions)
        selected_action_ids[str(p)] = selected_action_id
        selected_action_names.append(f"Player {p}: {allowed_actions[str(selected_action_id)]}")

    print(f"Will be sending: {selected_action_names}")
    if WAIT_FOR_KEY:
        input("Press enter to execute next actions")

    api.send_actions(id, selected_action_ids)
