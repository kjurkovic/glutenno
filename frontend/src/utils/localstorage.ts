export const LS_ACCESS_TOKEN = "gtAT"
export const LS_REFRESH_TOKEN = "gtRT"
export const LS_USER_ID = "gtID"

export const setAuthTokens = (accessToken: string, refreshToken: string) => {
  localStorage.setItem(LS_ACCESS_TOKEN, accessToken);
  localStorage.setItem(LS_REFRESH_TOKEN, refreshToken);
};

export const setUserId = (userId: string) => {
  localStorage.setItem(LS_USER_ID, userId);
};

export const getAccessToken = () => {
  return localStorage.getItem(LS_ACCESS_TOKEN);
};

export const getUserId = () => {
  return localStorage.getItem(LS_USER_ID);
};

export const getRefreshToken = () => {
  return localStorage.getItem(LS_REFRESH_TOKEN);
};

export const removeTokens = () => {
  localStorage.removeItem(LS_ACCESS_TOKEN);
  localStorage.removeItem(LS_REFRESH_TOKEN);
  localStorage.removeItem(LS_USER_ID);
};
