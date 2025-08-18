import axios from "axios"
import { ParentForm } from "./ParentForm"
import { useState } from "react"
import { useNavigate } from "react-router-dom";


export function Auth() {

    const [name, setName] = useState('');
    const [password, setPassword] = useState('');

    const navigate = useNavigate();

    const loginHandler = async () => {
        axios.interceptors.request.use(config => {
            config.withCredentials = true;
            return config;
        });
        const result = await axios.post(`http://localhost:5000/api/auth/login`, { 
            name,
            password
        })
        .catch((error) => {
            console.log(error);
        })

        if(result?.status === 200) {
            navigate('/chat', { replace: true });
        }
    }

    const registrationHandler = async () => {
        axios.interceptors.request.use(config => {
            config.withCredentials = true;
            return config;
        });
        const result = await axios.post(`http://localhost:5000/api/auth/reg`, { 
            name,
            password
        })
        .catch((error) => {
            console.log(error);
        })

        if(result?.status === 200) {
            navigate('/chat', { replace: true });
        }
    }

    return (
        <ParentForm
            isDialog={false}
            isOpen={true}
            setIsOpen={()=>{}}
        >
            <p>Authorisation</p>
            <input onChange={(event)=> {
                setName(event.target.value);
            }}></input>
            <input type="password" onChange={(event)=> {
                setPassword(event.target.value);
            }}></input>
            <br/>
            <button onClick={loginHandler}>Sign in</button>
            <button onClick={registrationHandler}>Sign up</button>

        </ParentForm>
    )
}

export default Auth