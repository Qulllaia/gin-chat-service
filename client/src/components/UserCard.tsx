import { User } from "../types";
import '../styles/UserCard.css'
interface UserCardProps {
    user: User
    createChatHandler: (user_id: number) =>  Promise<void>
}

export function UserCard({user, createChatHandler}: UserCardProps) {
    return (
        <div key={user.id} className="user-card"
            onClick={()=>createChatHandler(user.id)}
        >
            <p className="user-p">{user.name}</p>
        </div>
    )
}

export default UserCard;