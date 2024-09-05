special_roles = {
    "Programming Languages": {
        "🐍": "python",
        "🌐": "html",
        "☕": "java",
        "⚡": "zig",
        "🦦": "golang",
        "🦀": "rust",
        "🤖": "c",
        "💻": "c++",
        "🔷": "c#",
        "🎨": "css",
        "🟨": "javascript",
        "🌍": "web dev",
        "🛠️": "backend dev",
        "🧑‍💻": "fullstack dev",
        "⚙️": "systems engineering",
    },
    "Operating Systems": {
        "🪶": "macOS",
        "🪟": "windows",
        "🐧": "linux",
    },
    "Area of Interest": {
        "🌍": "web dev",
        "🛠️": "backend dev",
        "🧑‍💻": "fullstack dev",
        "⚙️": "systems engineering",
    },
}
all_roles_dict = {}

for key in special_roles:
    for role in special_roles[key]:
        all_roles_dict[role] = special_roles[key][role]

all_roles_list = ()

for key in special_roles:
    for role in special_roles[key]:
        all_roles_list += (role,)
