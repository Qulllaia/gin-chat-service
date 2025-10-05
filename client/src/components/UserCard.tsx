import { User } from "../types";
import '../styles/UserCard.css'
interface UserCardProps {
    user: User
    createChatHandler: (user_id: number) =>  Promise<void>
    handleCheckboxChange: (value: string) => void,
    checkboxAvaible: boolean,
}

export function UserCard({user, createChatHandler, handleCheckboxChange, checkboxAvaible}: UserCardProps) {
    return (
        <div key={user.id} className="user-card"
            onClick={()=> checkboxAvaible ? null : createChatHandler(user.id)}
        >
            <p className="user-p">{user.name}</p>
            {
                checkboxAvaible &&  
                <div className="checkbox">  
                    <input 
                        className="form-check-input m-3" 
                        type="checkbox" 
                        value="" 
                        id="flexCheckDefault" 
                        onChange={
                            ()=> checkboxAvaible ? handleCheckboxChange(user.id.toString()) : null
                        }
                    ></input>
                    <label className="form-check-label" htmlFor="flexCheckDefault"></label>
                </div>
            }
        </div>
    )
}

export default UserCard;