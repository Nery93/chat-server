import "./OnlineUsers.css";

interface OnlineUsersProps {
  users: string[];
  ownUsername: string;
}

export function OnlineUsers({ users, ownUsername }: OnlineUsersProps) {
  if (users.length === 0) return null;

  return (
    <div className="online-users">
      <span className="online-users__label">online</span>
      <ul className="online-users__list">
        {users.map((user) => (
          <li
            key={user}
            className={`online-users__pill ${user === ownUsername ? "online-users__pill--mine" : ""}`}
          >
            {user === ownUsername ? "você" : user}
          </li>
        ))}
      </ul>
    </div>
  );
}
