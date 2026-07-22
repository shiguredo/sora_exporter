# canary.py の git commit に -- VERSION を追加する

- Priority: Low
- Created: 2026-07-23
- Completed: {YYYY-MM-DD}
- Model: Qwen Code
- Branch: feature/fix-canary-py-git-commit
- Polished: {YYYY-MM-DD}

## 目的

canary.py の git commit がステージ済みの無関係な変更を巻き込まないようにする。

## 優先度根拠

`git add VERSION` の前に他のファイルがステージされている場合、バージョン上げコミットに無関係な変更が混入する。

## 現状

`canary.py:52-54`:

```python
subprocess.run(["git", "add", "VERSION"], check=True)
subprocess.run(
    ["git", "commit", "-m", f"[canary] バージョンを {new_version} にあげる"], check=True
)
```

`git commit -m "..."` はステージ領域の全変更をコミットする。

## 完了条件

- `git commit -m "..." -- VERSION` に修正されている

## 解決方法

`canary.py:53-54` の `git commit` コマンドに `-- VERSION` を追加する。
