import { match } from 'assert';
import React, { useEffect, useState } from 'react';
import { Link, NavLink, Route, Routes, Navigate, Outlet } from 'react-router-dom';
import { getAccessToken, removeTokens } from '../../utils/localstorage';
import Admin from '../pages/admin';
import Home from '../pages/home';
import Login from '../pages/login';
import Recipe from '../pages/recipe';
import RecipeNew from '../pages/recipenew';
import Recipes from '../pages/recipes';
import Register from '../pages/register';

interface IProtected {
  isLoggedIn: boolean,
  children: any,
}

const Protected = (props: IProtected) => {
  if (!props.isLoggedIn) {
    return <Outlet />
  }
  return props.children;
};

function App() {
  const [isLoggedIn, setLoggedIn] = useState(false)

  const handleSignout = () => {
    removeTokens()
    window.location.reload()
  }

  useEffect(() => {
    const token = getAccessToken() 
    setLoggedIn(token != null)
  })

  return (
    <div className="p-10">
        <div className="flex flex-row flex-wrap">
          <div className="grow">
            <Link className="text-2xl antialiased font-bold" to="/">Let's Cook!</Link>
          </div>
          <nav>
            <NavLink className="mr-6" to="/">Home</NavLink> 
            <NavLink  className="mr-6" to="/recipes">Recipes</NavLink> 
            <Protected isLoggedIn={isLoggedIn}>
              <NavLink className="mr-6" to="/admin">Admin</NavLink>
            </Protected>
            {
              isLoggedIn 
              ? <button className="mr-6 hover:underline" onClick={handleSignout}>Sign out</button> 
              : <><NavLink className="mr-6" to="/signin">Sign in</NavLink><NavLink className="mr-6" to="/register">Sign up</NavLink></> 
              }
          </nav>
        </div>
        <main className="mt-10">
          <Routes>
            <Route path="/" element={<Home/>} />
            <Route path="/signin" element={<Login/>} />
            <Route path="/register" element={<Register/>} />
            <Route path="/recipes" element={<Recipes />} />
            <Route path="/recipes/:recipeId" element={<Recipe />} />
            <Route path="/admin" element={
              <Protected isLoggedIn={isLoggedIn}>
                <Admin />
              </Protected>
            } />
            <Route path="/admin/recipe/add" element={
              <Protected isLoggedIn={isLoggedIn}>
                <RecipeNew />
              </Protected>
            } />
            <Route path="/admin/recipe/:recipeId" element={
              <Protected isLoggedIn={isLoggedIn}>
                <RecipeNew />
              </Protected>
            } />
          </Routes>
        </main>
    </div>
  );
}

export default App;
