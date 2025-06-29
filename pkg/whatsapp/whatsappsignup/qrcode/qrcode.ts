import QRCode from "qrcode";
import { getServerProps } from "../../../common/utils";

const { phoneNumber, message } = getServerProps<{ phoneNumber: string, message: string }>("whatsappQrCodeProps");

console.log({ phoneNumber, message });
const attachedMessage = message ? `?text=${encodeURIComponent(message)}` : ""
const whatsappUrl = `http://wa.me/${phoneNumber}${attachedMessage}`;

const qrCodeElement = document.getElementById("whatsapp-qr-code") as HTMLCanvasElement;
if (qrCodeElement) {
	QRCode.toCanvas(qrCodeElement, whatsappUrl, {
		width: 288,
		margin: 2,
		color: {
			dark: '#000000',
			light: '#FFFFFF'
		},
	}).catch((error) => {
		console.error('Error generating QR code:', error);
	})
} else {
	console.error("QR code canvas element not found");
}

