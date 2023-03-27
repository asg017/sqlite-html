from datasette.app import Datasette
import pytest


@pytest.mark.asyncio
async def test_plugin_is_installed():
    datasette = Datasette(memory=True)
    response = await datasette.client.get("/-/plugins.json")
    assert response.status_code == 200
    installed_plugins = {p["name"] for p in response.json()}
    assert "datasette-sqlite-html" in installed_plugins

@pytest.mark.asyncio
async def test_sqlite_html_functions():
    datasette = Datasette(memory=True)
    response = await datasette.client.get("/_memory.json?sql=select+html_version()")
    assert response.status_code == 200
    html_version, = response.json()["rows"][0]
    assert html_version[0] == "v"