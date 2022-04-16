import asyncio
import time

import pytest
import requests
from telethon import TelegramClient
from telethon.tl.custom import Conversation
from telethon.tl.custom.message import Message

baseUrl = "http://localhost:8090"
api_id = 11519316
api_hash = '425a717993b73e9671e2b254a6461e31'
bot = "@SpbstuDatingBot"
tgSession = "test"
userId = "l0gark"

start_answer = "Привет! Я, SpbstuDatingBot, помогаю людям познакомиться\n\n" \
               + "Список доступных команд: \n" \
               + "- /start - начало работы\n" \
               + "- /profile - заполнить анкету\n" \
               + "- /next - показать следующего пользователя"


def clear_system():
    requests.post(baseUrl + "/deleteAll")


async def sendStart(conv: Conversation):
    await conv.send_message("/start")
    resp: Message = await conv.get_response()
    assert start_answer == resp.raw_text


async def sendProfile(conv: Conversation, sex):
    fullSex = "Женщина"
    if sex == "М":
        fullSex = "Мужчина"

    await conv.send_message("/profile")
    resp: Message = await conv.get_response()
    assert resp.raw_text == "Как Вас зовут?"
    await conv.send_message("name")
    resp: Message = await conv.get_response()
    assert resp.raw_text == "Сколько Вам лет?"
    await conv.send_message("1")
    resp: Message = await conv.get_response()
    assert resp.raw_text == "Из какого Вы города?"
    await conv.send_message("City")
    resp: Message = await conv.get_response()
    assert resp.raw_text == "Введите краткое описание своего профиля."
    await conv.send_message("Description")
    resp: Message = await conv.get_response()
    assert resp.raw_text == "Пришлите фотографию, которая будет показываться другим пользователям в ленте."
    await conv.send_file("img.jpg")
    resp: Message = await conv.get_response()
    assert resp.raw_text == "Какого Вы пола? М/Ж"
    await conv.send_message(sex)
    resp: Message = await conv.get_response()
    assert resp.raw_text == "Имя: name\nВозраст: 1\nГород: City\nОписание: Description\nПол: " + fullSex + "\n\nПопробуйте ввести команду /next"
    return resp.raw_text


@pytest.fixture(scope="session")
def event_loop():
    return asyncio.get_event_loop()


@pytest.fixture(scope="session")
async def client() -> TelegramClient:
    client = TelegramClient(
        tgSession, api_id, api_hash,
        sequential_updates=True
    )
    # Connect to the server
    await client.connect()
    # Issue a high level command to start receiving message
    await client.get_me()
    # Fill the entity cache
    await client.get_dialogs()

    yield client

    await client.disconnect()
    await client.disconnected


@pytest.mark.asyncio
async def test_scenario1_1(client: TelegramClient):
    clear_system()

    # Create a conversation
    async with client.conversation(bot, timeout=5) as conv:
        # Send a command
        await conv.send_message("/next")
        resp: Message = await conv.get_response()
        exp = start_answer
        assert exp == resp.raw_text


@pytest.mark.asyncio
async def test_scenario1_2(client: TelegramClient):
    clear_system()
    # Create a conversation
    async with client.conversation(bot, timeout=5) as conv:
        # Send a command
        await conv.send_message("/profile")
        resp: Message = await conv.get_response()
        exp = start_answer
        assert exp == resp.raw_text


@pytest.mark.asyncio
async def test_scenario2(client: TelegramClient):
    clear_system()
    # Create a conversation
    async with client.conversation(bot, timeout=5) as conv:
        # Send a command
        await conv.send_message("/start")
        resp: Message = await conv.get_response()
        exp = start_answer
        assert exp == resp.raw_text


@pytest.mark.asyncio
async def test_scenario3(client: TelegramClient):
    clear_system()
    # Create a conversation
    async with client.conversation(bot, timeout=5) as conv:
        # Send a command
        await sendStart(conv)
        await conv.send_message("/start")
        resp: Message = await conv.get_response()
        exp = "Вы уже зарегистрированы в системе"
        assert exp == resp.raw_text


@pytest.mark.asyncio
async def test_scenario4(client: TelegramClient):
    clear_system()
    # Create a conversation
    async with client.conversation(bot, timeout=5) as conv:
        # Send a command
        await sendStart(conv)
        await conv.send_message("/profile")
        resp: Message = await conv.get_response()
        exp = "Как Вас зовут?"
        assert exp == resp.raw_text


@pytest.mark.asyncio
async def test_scenario5(client: TelegramClient):
    clear_system()
    # Create a conversation
    async with client.conversation(bot, timeout=5) as conv:
        await sendStart(conv)
        await sendProfile(conv, "М")
        await conv.send_message("/profile")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Как Вас зовут?"


@pytest.mark.asyncio
async def test_scenario6(client: TelegramClient):
    clear_system()
    # Create a conversation
    async with client.conversation(bot, timeout=5) as conv:
        # Send a command
        await sendStart(conv)
        await conv.send_message("/profile")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Как Вас зовут?"
        await conv.send_file('file.txt')
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Данные введены некорректно, попробуйте снова."


# False - Ж, True - М
@pytest.mark.asyncio
async def test_scenario7(client: TelegramClient):
    clear_system()
    requests.post(baseUrl + "/addTestUser?sex=false")
    # Create a conversation
    async with client.conversation(bot, timeout=5) as conv:
        await sendStart(conv)
        await sendProfile(conv, "М")
        await conv.send_message("/next")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Имя: TestName\nВозраст: 0\nГород: TestCity\nОписание: TestDescription\nПол: Женщина"


@pytest.mark.asyncio
async def test_scenario8(client: TelegramClient):
    clear_system()
    # Create a conversation
    async with client.conversation(bot, timeout=5) as conv:
        await sendStart(conv)
        await sendProfile(conv, "М")
        await conv.send_message("/next")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Все анкеты просмотрены. Попробуйте ещё раз немного позже."


@pytest.mark.asyncio
async def test_scenario10(client: TelegramClient):
    clear_system()
    requests.post(baseUrl + "/addTestUser?sex=false")
    # Create a conversation
    async with client.conversation(bot, timeout=5) as conv:
        await sendStart(conv)
        await sendProfile(conv, "М")
        await conv.send_message("/next")
        resp: Message = await conv.get_response()
        assert resp.raw_text != "Все анкеты просмотрены. Попробуйте ещё раз немного позже."
        await resp.click(0, 1)
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Все анкеты просмотрены. Попробуйте ещё раз немного позже."


@pytest.mark.asyncio
async def test_scenario11(client: TelegramClient):
    clear_system()
    requests.post(baseUrl + "/addTestUser?sex=false")
    # Create a conversation
    async with client.conversation(bot, timeout=5) as conv:
        await sendStart(conv)
        await sendProfile(conv, "М")
        await conv.send_message("/next")
        resp: Message = await conv.get_response()
        assert resp.raw_text != "Все анкеты просмотрены. Попробуйте ещё раз немного позже."
        await resp.click(0, 0)
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Все анкеты просмотрены. Попробуйте ещё раз немного позже."


@pytest.mark.asyncio
async def test_scenario13(client: TelegramClient):
    clear_system()
    # Create a conversation
    async with client.conversation(bot, timeout=5) as conv:
        await sendStart(conv)
        await conv.send_message("/profile")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Как Вас зовут?"
        await conv.send_message("name")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Сколько Вам лет?"
        await conv.send_message("1")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Из какого Вы города?"
        await conv.send_message("City")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Введите краткое описание своего профиля."
        await conv.send_message("Description")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Пришлите фотографию, которая будет показываться другим пользователям в ленте."
        await conv.send_message("not photo")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Данные введены некорректно, попробуйте снова."


@pytest.mark.asyncio
async def test_scenario15(client: TelegramClient):
    clear_system()
    # Create a conversation
    async with client.conversation(bot, timeout=5) as conv:
        await sendStart(conv)
        await conv.send_message("/wrong")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Такой команды не существует.\n\nСписок доступных команд: \n- /start - начало работы\n- /profile - заполнить анкету\n- /next - показать следующего пользователя"


@pytest.mark.asyncio
async def test_scenario17(client: TelegramClient):
    clear_system()
    # Create a conversation
    async with client.conversation(bot, timeout=5) as conv:
        await sendStart(conv)
        await sendProfile(conv, "Ж")
        await sendProfile(conv, "М")


@pytest.mark.asyncio
async def test_scenario18(client: TelegramClient):
    clear_system()
    # Create a conversation
    async with client.conversation(bot, timeout=5) as conv:
        await sendStart(conv)
        await sendProfile(conv, "Ж")


@pytest.mark.asyncio
async def test_scenario19(client: TelegramClient):
    clear_system()
    # Create a conversation
    async with client.conversation(bot, timeout=5) as conv:
        await sendStart(conv)
        exp = await sendProfile(conv, "Ж")
        await conv.send_message("/profile")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Как Вас зовут?"
        await resp.click(0)
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Сколько Вам лет?"
        await resp.click(0)
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Из какого Вы города?"
        await resp.click(0)
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Введите краткое описание своего профиля."
        await resp.click(0)
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Пришлите фотографию, которая будет показываться другим пользователям в ленте."
        await resp.click(0)
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Какого Вы пола? М/Ж"
        await resp.click(0)
        resp: Message = await conv.get_response()
        assert exp == resp.raw_text


@pytest.mark.asyncio
async def test_scenario20(client: TelegramClient):
    clear_system()
    requests.post(baseUrl + "/addTestUser?sex=false")
    # Create a conversation
    async with client.conversation(bot, timeout=5) as conv:
        await sendStart(conv)
        await sendProfile(conv, "Ж")
        await conv.send_message("/next")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Все анкеты просмотрены. Попробуйте ещё раз немного позже."


@pytest.mark.asyncio
async def test_scenario22(client: TelegramClient):
    clear_system()
    # Create a conversation
    async with client.conversation(bot, timeout=5) as conv:
        await sendStart(conv)
        await conv.send_message("/profile")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Как Вас зовут?"
        await conv.send_message("name")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Сколько Вам лет?"
        await conv.send_message("1")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Из какого Вы города?"
        await conv.send_message("City")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Введите краткое описание своего профиля."
        await conv.send_message("Description")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Пришлите фотографию, которая будет показываться другим пользователям в ленте."
        await conv.send_file("img.jpg")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Какого Вы пола? М/Ж"
        await conv.send_message("WRONG")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Данные введены некорректно, попробуйте снова."


@pytest.mark.asyncio
async def test_scenario23(client: TelegramClient):
    clear_system()
    # Create a conversation
    async with client.conversation(bot, timeout=5) as conv:
        await sendStart(conv)
        await conv.send_message("/profile")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Как Вас зовут?"
        await conv.send_message("/next")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Пожалуйста дозаполните анкету."


@pytest.mark.asyncio
async def test_scenario25(client: TelegramClient):
    clear_system()
    # Create a conversation
    async with client.conversation(bot, timeout=5) as conv:
        await sendStart(conv)
        await sendProfile(conv, "М")
        await conv.send_message("/profile")
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Как Вас зовут?"
        await resp.click(0)
        resp: Message = await conv.get_response()
        assert resp.raw_text == "Сколько Вам лет?"


def test_regeneration():
    n = 10
    times = [0.0] * n

    for i in range(n):
        try:
            requests.post('http://localhost:8090/panic', verify=False, timeout=1)
        except requests.exceptions.ConnectionError as e:
            pass

        start = time.time()

        success = False
        max_recover_time = 10.0

        while not success or time.time() - start >= max_recover_time:
            try:
                r = requests.get('http://localhost:8090/ping')
                r.raise_for_status()
                success = True
            except requests.exceptions.ConnectionError:
                pass

        end = time.time()
        recover_time = end - start
        assert recover_time < max_recover_time
        times[i] = recover_time

        print('\n' + str(i) + '.', recover_time)

    print('\nAverage recovering time =', sum(times) / n)


def test_regeneration():
    n = 3
    times = [0.0] * n

    for i in range(n):
        try:
            requests.post('http://localhost:8090/panic')
        except requests.exceptions.ConnectionError as e:
            pass

        start = time.time()

        success = False
        max_recover_time = 10.0

        while not success or time.time() - start >= max_recover_time:
            try:
                r = requests.get('http://localhost:8090/ping')
                r.raise_for_status()
                success = True
            except requests.exceptions.ConnectionError:
                pass

        end = time.time()
        recover_time = end - start
        assert recover_time < max_recover_time
        times[i] = recover_time

        print('\n' + str(i) + '.', recover_time)

    print('\nAverage recovering time =', sum(times) / n)
