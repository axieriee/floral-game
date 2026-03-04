"""Terminal UI helpers for the text RPG."""

import os


def clear_screen():
    os.system("cls" if os.name == "nt" else "clear")


class TerminalUI:
    """Simple terminal-based UI for the game."""

    def __init__(self):
        self.width = 60

    def print_line(self, text=""):
        print(text)

    def print_banner(self, text):
        print()
        print("=" * self.width)
        padding = (self.width - len(text)) // 2
        print(" " * padding + text)
        print("=" * self.width)
        print()

    def print_box(self, lines):
        print("+" + "-" * (self.width - 2) + "+")
        for line in lines:
            trimmed = line[:self.width - 4]
            print(f"| {trimmed:<{self.width - 4}} |")
        print("+" + "-" * (self.width - 2) + "+")

    def print_separator(self):
        print("-" * self.width)

    def get_input(self, prompt="> "):
        try:
            return input(prompt).strip()
        except (EOFError, KeyboardInterrupt):
            print()
            return "quit"

    def get_choice(self, max_val, prompt="> "):
        while True:
            raw = self.get_input(prompt)
            if raw.lower() in ("q", "quit", "exit"):
                return -1
            try:
                val = int(raw)
                if 1 <= val <= max_val:
                    return val
            except ValueError:
                pass
            self.print_line(f"Please enter a number between 1 and {max_val}.")

    def confirm(self, question):
        self.print_line(f"{question} (y/n)")
        resp = self.get_input()
        return resp.lower() in ("y", "yes")
