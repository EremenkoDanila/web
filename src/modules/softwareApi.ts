export interface ISoftware {
    software_id: number;
    deleted_flg: boolean;
    title: string;
    description: string;
    full_description: string;
    img_url: string;
    os: string;
    size: number;
    version: string;
    functions: string;
}

export const getSoftwareByName = async (name: string = ""): Promise<ISoftware[]> => {
    try {
        // Используем относительный путь - Vite proxy обработает его
        const url = `/api/software?app=${encodeURIComponent(name)}`;
        console.log("Fetching from:", url);
        
        const response = await fetch(url);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        console.log("API data received:", data);
        return data;
    } catch (error) {
        console.error('Error fetching software:', error);
        throw error;
    }
};

export const getSoftwareById = async (id: number): Promise<ISoftware> => {
    try {
        const response = await fetch(`/api/software/${id}`);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error fetching software by id:', error);
        throw error;
    }
};