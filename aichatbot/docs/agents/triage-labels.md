# Triage Labels

The skills speak in terms of five canonical triage roles. This file maps those roles to the actual label strings used in this repo's issue tracker.

| Label in mattpocock/skills | Label in our tracker | Meaning                                  |
| -------------------------- | -------------------- | ---------------------------------------- |
| `needs-triage`             | `needs-triage`       | Maintainer needs to evaluate this issue  |
| `needs-info`               | `needs-info`         | Waiting on reporter for more information |
| `ready-for-agent`          | `ready-for-agent`    | Fully specified, ready for an AFK agent  |
| `ready-for-human`          | `ready-for-human`    | Requires human implementation            |
| `wontfix`                  | `wontfix`            | Will not be actioned                     |

When a skill mentions a role (e.g. "apply the AFK-ready triage label"), use the corresponding label string from this table.

Edit the right-hand column to match whatever vocabulary you actually use.

## 與 issue→MR 迴圈共用

這五個 label 也是 `rory-ai` 迴圈（`~/Workspace/rory/gitlab-issue-loop`）用的 triage 狀態，同一張 issue 同時只留一個。迴圈只處理指派給 `rory-ai` 的 issue，釐清後會改成 `needs-info`、`ready-for-agent` 或 `ready-for-human`；標了 `ready-for-human`、`wontfix` 的不會接。執行進度另外用 `ai::doing`、`ai::review`、`ai::blocked` 表示。
