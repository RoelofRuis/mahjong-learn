import random


def pick_first_action(available_actions: dict) -> int:
    return 0


def pick_random_action(available_actions: dict) -> int:
    return random.choice(list(available_actions.keys()))
