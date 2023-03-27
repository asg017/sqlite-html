from setuptools import setup

version = {}
with open("datasette_sqlite_html/version.py") as fp:
    exec(fp.read(), version)

VERSION = version['__version__']

setup(
    name="datasette-sqlite-html",
    description="",
    long_description="",
    long_description_content_type="text/markdown",
    author="Alex Garcia",
    url="https://github.com/asg017/sqlite-html",
    project_urls={
        "Issues": "https://github.com/asg017/sqlite-html/issues",
        "CI": "https://github.com/asg017/sqlite-html/actions",
        "Changelog": "https://github.com/asg017/sqlite-html/releases",
    },
    license="MIT License, Apache License, Version 2.0",
    version=VERSION,
    packages=["datasette_sqlite_html"],
    entry_points={"datasette": ["sqlite_html = datasette_sqlite_html"]},
    install_requires=["datasette", "sqlite-html"],
    extras_require={"test": ["pytest"]},
    python_requires=">=3.7",
)