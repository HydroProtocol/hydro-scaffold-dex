/**
 * Login Data
 * {
 *   address: "0x....",
 *   hydroAuthentication: "xxx.bbb.ccc"
 * }
 */

export const saveLoginData = (address, hydroAuthentication) => {
  window.localStorage.setItem(`loginData-${address}`, JSON.stringify({ address, hydroAuthentication }));
};

export const cleanLoginDate = address => {
  window.localStorage.removeItem(`loginData-${address}`);
};

export const loadAccountHydroAuthentication = address => {
  const savedData = window.localStorage.getItem(`loginData-${address}`);

  if (!savedData) {
    return null;
  }

  let loginData;
  try {
    loginData = JSON.parse(savedData);
  } catch (e) {
    cleanLoginDate(address);
    return null;
  }

  if (loginData.address && loginData.address.toLowerCase() === address.toLowerCase()) {
    return loginData.hydroAuthentication;
  }

  return null;
};
