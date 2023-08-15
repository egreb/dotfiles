import os
import re
from typing import Match, Optional, Pattern

from ..mark import Mark
from ..targets.file_target import FileTarget, ContentType
from .finder import BaseFinder


class RailsLogPartialFinder(BaseFinder):
    """
    From a text that looks like this:
        Rendered layouts/_base.html.erb (Duration: 32.9ms | Allocations: 2204)

    this finder extracts file path of a rails controller (with line number pointing to action method).
    """

    @classmethod
    def pattern(cls) -> Pattern[str]:
        return re.compile(
            r'Render(?:ed|ing) ([-a-zA-Z0-9_+-,./]+)'
        )

    def match_to_mark(self, match: Match[str]) -> Optional[Mark]:
        start = match.span(1)[0]
        text = match.group(1)
        file_path = os.path.join(self.path_prefix, 'app/views/' + text)

        if os.path.exists(file_path):
            return Mark(
                start=start,
                text=text,
                target=FileTarget(file_path, ContentType.TEXT)
            )

        return None
