// services/ProTeeService.js
export class ProTeeService {
    constructor(apiClient, eventBus) {
        this.api = apiClient;
        this.eventBus = eventBus;
        this.status = null;
    }

    async start(path) {
        try {
            const body = path ? { path } : {};
            const response = await this.api.post('/api/protee/start', body);

            if (response.ok) {
                this.eventBus.emit('protee:starting');
                return { success: true };
            } else {
                const text = await response.text();
                throw new Error(text || `Failed to start: ${response.statusText}`);
            }
        } catch (error) {
            this.eventBus.emit('protee:error', error.message);
            return { success: false, error: error.message };
        }
    }

    async stop() {
        try {
            const response = await this.api.post('/api/protee/stop');

            if (response.ok) {
                this.eventBus.emit('protee:stopping');
                return { success: true };
            }
        } catch (error) {
            this.eventBus.emit('protee:error', error.message);
            return { success: false, error: error.message };
        }
    }

    async saveConfig(enabled, shotsPath) {
        try {
            const response = await this.api.post('/api/protee/config', {
                enabled, shotsPath
            });

            if (!response.ok) {
                throw new Error(`Failed to save config: ${response.statusText}`);
            }

            return { success: true };
        } catch (error) {
            this.eventBus.emit('protee:error', error.message);
            return { success: false, error: error.message };
        }
    }

    updateStatus(status) {
        this.status = status;
        this.eventBus.emit('protee:status', status);
    }
}
