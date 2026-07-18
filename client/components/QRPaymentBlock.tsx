"use client";

import { QRCodeSVG } from "qrcode.react";
import { CheckCircle, Clock, RefreshCw, Download, Camera } from "lucide-react";
import { useState, useEffect, useCallback } from "react";

interface QRPaymentData {
  order_id: string;
  qr_code_url?: string; // hosted image (may be empty in sandbox)
  qr_code_string?: string; // raw QRIS EMV string → rendered as QR
  payment_url?: string; // Duitku hosted payment page
  expires_at?: string;
  amount: number;
  gateway: string;
}

interface OrderItem {
  id: string;
  photo: { filename: string };
  download_url?: string;
}

interface OrderStatus {
  payment_status: "pending" | "paid" | "failed" | "expired";
  items?: OrderItem[];
  download_expires_at?: string;
}

interface QRPaymentBlockProps {
  qr: QRPaymentData;
  buyerEmail: string;
  onPaid?: () => void;
}

const API_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api/v1";

export default function QRPaymentBlock({
  qr,
  buyerEmail,
  onPaid,
}: QRPaymentBlockProps) {
  const [status, setStatus] = useState<
    "waiting" | "paid" | "failed" | "expired"
  >("waiting");
  const [orderStatus, setOrderStatus] = useState<OrderStatus | null>(null);
  const [lastCheck, setLastCheck] = useState<Date | null>(null);

  const formatRupiah = (n: number) =>
    new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(n);

  const checkStatus = useCallback(async () => {
    try {
      const res = await fetch(
        `${API_URL}/public/orders/${qr.order_id}?email=${encodeURIComponent(buyerEmail)}`,
      );
      const data = await res.json();
      const order: OrderStatus = data.data;
      setOrderStatus(order);
      setLastCheck(new Date());

      if (order.payment_status === "paid") {
        setStatus("paid");
        onPaid?.();
      } else if (order.payment_status === "failed") {
        setStatus("failed");
      } else if (order.payment_status === "expired") {
        setStatus("expired");
      }
    } catch {
      // Silently ignore network errors during polling
    }
  }, [qr.order_id, buyerEmail, onPaid]);

  // Poll every 5 seconds while waiting
  useEffect(() => {
    if (status !== "waiting") return;
    const id = setInterval(checkStatus, 5000);
    return () => clearInterval(id);
  }, [status, checkStatus]);

  if (status === "paid") {
    return (
      <div className="text-center py-8">
        <CheckCircle size={56} className="text-green-500 mx-auto mb-4" />
        <h2 className="text-2xl font-bold text-gray-900 mb-2">
          Pembayaran Berhasil! 🎉
        </h2>
        <p className="text-gray-600 mb-6">
          Link download telah dikirim ke <strong>{buyerEmail}</strong>
        </p>

        {orderStatus?.items && orderStatus.items.length > 0 && (
          <div className="bg-gray-50 rounded-xl p-4 mb-6 text-left space-y-2 max-w-md mx-auto">
            <p className="text-sm font-semibold text-gray-700 mb-3">
              Download Foto Anda:
            </p>
            {orderStatus.items.map((item) =>
              item.download_url ? (
                <a
                  key={item.id}
                  href={item.download_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center gap-2 text-violet-600 hover:text-violet-800 text-sm py-1"
                >
                  <Download size={14} />
                  <span className="truncate">{item.photo?.filename}</span>
                </a>
              ) : null,
            )}
            {orderStatus.download_expires_at && (
              <p className="text-xs text-gray-400 mt-2">
                Link aktif hingga:{" "}
                {new Date(orderStatus.download_expires_at).toLocaleDateString(
                  "id-ID",
                  {
                    day: "numeric",
                    month: "long",
                    year: "numeric",
                  },
                )}
              </p>
            )}
          </div>
        )}
      </div>
    );
  }

  if (status === "expired") {
    return (
      <div className="text-center py-8">
        <Clock size={48} className="text-amber-400 mx-auto mb-4" />
        <h2 className="text-xl font-bold text-gray-900 mb-2">
          QR Code Kedaluwarsa
        </h2>
        <p className="text-gray-600 mb-4">
          Silakan buat pesanan baru untuk melanjutkan pembayaran.
        </p>
        <button
          onClick={() => window.location.reload()}
          className="px-6 py-2 bg-violet-600 text-white rounded-xl font-medium hover:bg-violet-700"
        >
          Buat Pesanan Baru
        </button>
      </div>
    );
  }

  return (
    <div className="text-center">
      {/* Amount */}
      <div className="bg-violet-50 rounded-2xl p-4 mb-6">
        <p className="text-sm text-violet-600 font-medium">Total Pembayaran</p>
        <p className="text-3xl font-bold text-violet-700 mt-1">
          {formatRupiah(qr.amount)}
        </p>
        <p className="text-xs text-violet-500 mt-1">via QRIS · {qr.gateway}</p>
      </div>

      {/* QR Code */}
      <div className="flex justify-center mb-4">
        <div className="bg-white rounded-2xl p-4 shadow-lg border-2 border-violet-100 inline-block">
          {qr.qr_code_string ? (
            <QRCodeSVG
              value={qr.qr_code_string}
              size={220}
              level="M"
              includeMargin
              imageSettings={{
                src: "/logo.png",
                height: 32,
                width: 32,
                excavate: true,
              }}
            />
          ) : qr.qr_code_url ? (
            /* Fallback: hosted image from Midtrans */
            // eslint-disable-next-line @next/next/no-img-element
            <img
              src={qr.qr_code_url}
              alt="QRIS QR Code"
              className="w-52 h-52 object-contain"
            />
          ) : (
            <div className="w-52 h-52 flex items-center justify-center bg-gray-100 rounded-xl">
              <div className="text-center text-gray-400">
                <Camera size={32} className="mx-auto mb-2" />
                <p className="text-xs">QR tidak tersedia</p>
                <p className="text-xs">(sandbox mode)</p>
              </div>
            </div>
          )}
        </div>
      </div>

      {qr.payment_url && (
        <a
          href={qr.payment_url}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center justify-center rounded-full bg-violet-600 px-4 py-2 text-sm font-semibold text-white hover:bg-violet-700 mb-4"
        >
          Buka halaman pembayaran
        </a>
      )}

      {/* Instructions */}
      <div className="text-sm text-gray-600 mb-4 space-y-1">
        <p className="font-medium text-gray-800">Cara Bayar:</p>
        <p>1. Buka aplikasi perbankan/dompet digital Anda</p>
        <p>2. Pilih fitur Scan QR / QRIS</p>
        <p>3. Scan QR code di atas</p>
        <p>4. Konfirmasi pembayaran</p>
      </div>

      {/* Supported apps */}
      <div className="flex flex-wrap justify-center gap-2 text-xs text-gray-400 mb-6">
        {[
          "GoPay",
          "OVO",
          "Dana",
          "ShopeePay",
          "LinkAja",
          "BCA",
          "BRI",
          "BNI",
          "Mandiri",
        ].map((app) => (
          <span key={app} className="px-2 py-0.5 bg-gray-100 rounded-full">
            {app}
          </span>
        ))}
      </div>

      {/* Polling status */}
      <div className="flex items-center justify-center gap-2 text-sm text-gray-500">
        <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse" />
        <span>Menunggu konfirmasi pembayaran…</span>
        <button
          onClick={checkStatus}
          title="Periksa status"
          className="text-violet-500 hover:text-violet-700"
        >
          <RefreshCw size={14} />
        </button>
      </div>
      {lastCheck && (
        <p className="text-xs text-gray-400 mt-1">
          Terakhir dicek: {lastCheck.toLocaleTimeString("id-ID")}
        </p>
      )}
    </div>
  );
}
