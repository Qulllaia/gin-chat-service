import axios from "axios"
import { ParentForm } from "./ParentForm"
import { useState } from "react"
import { useNavigate } from "react-router-dom";


export function Auth() {

    const [name, setName] = useState('');
    const [password, setPassword] = useState('');
    const [email, setEmail] = useState('');
    const [authError, setAuthError] = useState('');

    const navigate = useNavigate();

    const loginHandler = async () => {
        axios.interceptors.request.use(config => {
            config.withCredentials = true;
            return config;
        });
        const result = await axios.post(`http://localhost:5000/api/auth/login`, { 
            email, 
            name,
            password
        })
        .catch((error) => {
            if(error.status === 500) {
                setAuthError('Fatal Error')
            } else {
                setAuthError(error.response.data.message)
            }
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
        const result = await axios.post(`http://localhost:5000/api/auth/verify`, { 
            email, 
            name,
            password
        })
        .catch((error) => {
            if(error.status === 500) {
                setAuthError('Fatal Error')
            } else {
                setAuthError(error.response.data.message)
            } 
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
            <form> 
                <h1 className="h3 mb-3 fw-normal">Please sign in</h1> 
                <div className="form-floating"> 
                    <input className="form-control" id="floatingInput"
                        onChange={(event)=> {
                            setEmail(event.target.value);
                        }}/> 
                    <label htmlFor="floatingInput">Email</label> 
                </div> 

                <div className="form-floating"> 
                    <input className="form-control" id="floatingInput"
                        onChange={(event)=> {
                            setName(event.target.value);
                        }}/> 
                    <label htmlFor="floatingInput">Login</label> 
                </div> 
                
                <div className="form-floating"> 
                    <input type="password" className="form-control mb-2" id="floatingPassword" placeholder="Password"  
                        onChange={(event)=> {
                            setPassword(event.target.value);
                        }}/> 
                    <label htmlFor="floatingPassword">Password</label> 
                </div> 
                <button className="btn btn-primary w-100 py-2 mb-2" type="button" onClick={loginHandler}>Sign in</button>
                <button className="btn btn-primary w-100 py-2" type="button" onClick={registrationHandler}>Sign up</button> 
            </form>
            <h5 className="auth-error">{authError}</h5>
        </ParentForm>
    )
}

export default Auth