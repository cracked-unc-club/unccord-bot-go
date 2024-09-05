# invite link:
# https://discord.com/api/oauth2/authorize?client_id={client_id}&permissions=8&scope=bot
# run on Heroku server and locally

from email import message

import discord
from discord.ext import commands
from discord.utils import get

import os

from dotenv import load_dotenv
from constants import all_roles_dict, all_roles_list, special_roles

load_dotenv()

intents = discord.Intents.all()
client = commands.Bot(command_prefix=commands.when_mentioned_or("$"), intents=intents)


def getToken():
    # code to open and read token
    return os.environ.get("TOKEN")


@client.event
async def on_ready():
    await client.change_presence(
        status=discord.Status.online,
        activity=discord.Game(
            name="React to my messages in #react-roles to show how cracked you are"
        ),
    )
    print("My body is ready")
    print("We have logged in as {0.user}".format(client))
    print("Name: {}".format(client.user.name))
    print("ID: {}".format(client.user.id))


@client.event
async def on_message(message):
    await client.process_commands(message)
    if (message.author.bot) and (
        message.author.id == client.user.id
    ):  # checks if message is from bot
        print(message.embeds)
        for embed in message.embeds:
            for field in embed.fields:
                value = field.value
                for i in range(len(all_roles_list)):
                    if all_roles_list[i] in value:
                        await message.add_reaction(all_roles_list[i])
    else:
        print(message.author.name)


@client.event
async def on_raw_reaction_add(payload):
    if payload.user_id != client.user.id:  # checks if reaction is from bot
        print("Add role initiated")
        guild_id = payload.guild_id
        guild = discord.utils.find(lambda g: g.id == guild_id, client.guilds)
        value = payload.emoji.name
        if all_roles_dict[value] is not None:
            role = discord.utils.get(guild.roles, name=all_roles_dict[value])
            while role is None:
                await guild.create_role(name=all_roles_dict[value], mentionable=True)
                role = discord.utils.get(guild.roles, name=all_roles_dict[value])
                await role.edit(mentionable=True)
                print(role)
            member = payload.member
            if member:
                await member.add_roles(role)
                print("success")
            else:
                print("Member not found")
        else:
            print("Role not found")
    else:
        print("The bot reacted")
        print(payload.user_id)


@client.event
async def on_raw_reaction_remove(payload):
    if payload.user_id != client.user.id:  # checks if reaction is from bot
        print("Add role initiated")
        guild_id = payload.guild_id
        guild = discord.utils.find(lambda g: g.id == guild_id, client.guilds)
        value = payload.emoji.name
        if all_roles_dict[value] is not None:
            role = discord.utils.get(guild.roles, name=all_roles_dict[value])
            while role is None:
                await guild.create_role(name=all_roles_dict[value], mentionable=True)
                role = discord.utils.get(guild.roles, name=all_roles_dict[value])
                await role.edit(mentionable=True)
                print(role)
            member = get(guild.members, id=payload.user_id)
            if member:
                await member.remove_roles(role)
                print("success")
            else:
                print("Member not found")
        else:
            print("Role not found")
    else:
        print("The bot reacted")
        print(payload.user_id)


@client.command()
async def setroles(ctx):
    channel = ctx.message.channel
    if ctx.message.author.guild_permissions.administrator:
        async for message in channel.history(
            limit=10000
        ):  # clears all of the bot messages in the channel
            if message.author == client.user:
                await message.delete()
        await ctx.message.delete()
        global special_roles
        # TODO remove previous messages the bot sent in the channel
        for key in special_roles.keys():
            embed = discord.Embed(color=0xFF0000)
            embed.add_field(
                name=f"Role Menu: {key}",
                value="React to give yourself a role.",
                inline=False,
            )
            for subkeys in special_roles[key].keys():
                embed.add_field(
                    name="\u200b",
                    value=f"{subkeys} : {special_roles[key][subkeys]}",
                    inline=False,
                )
            await ctx.send(embed=embed)
    else:
        msg = f"Sorry {ctx.message.author.mention}, only Admins can use this command"
        await channel.send(msg)


@client.command()
async def test(ctx, *, arg):
    print(arg)
    await ctx.send(arg)


@client.command()
async def ping(ctx):
    print(ctx)
    await ctx.send(f"Pong! {round(client.latency * 1000)} ms")


client.run(getToken())
