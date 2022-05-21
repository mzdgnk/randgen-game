// state = {
//     room: {
//         id: "xxx-xxx-xxx",
//         started: false,
//         topic: "",
//         users: [
//             { name: "user01", num: 50, open: false, isYou: true },
//             { name: "user02", num: 30, open: false, isYou: false },
//         ],
//     },
//     you: "user01",
//     joined: true,
// };

const room = (state, action) => {
    let s = addIsYou(internalRoom(state, action));
    console.log({ s });
    return s;
};

const internalRoom = (state, action) => {
    console.log({ state });
    console.log({ action });
    switch (action.type) {
        case "UPDATE":
            return { ...state, room: action.room };
        case "ASSIGN_USER":
            return { ...state, you: action.you };
        default:
            return state;
    }
};

export default room;

const addIsYou = (state) => {
    let users = [];
    let joined = false;
    for (const user of state.room.users) {
        console.log({ user });
        if (user.name == state.you) {
            users.push({ ...user, isYou: true });
            joined = true;
        } else {
            users.push({ ...user, isYou: false });
        }
    }
    return {
        ...state,
        room: {
            ...state.room,
            users: users,
        },
        joined: joined,
    };
};
